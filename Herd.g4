grammar Herd;

// Keep this  the top, as it gobbles up everything after it
RUN: 'run' ~('\n')*;
SET: 'set' ;
ADD: 'add' ;
REMOVE: 'remove' ;
LIST: 'list' ;
HOSTS: 'hosts' ;
ONELINE: '--oneline' ;
DURATION: ( '-'? [0-9]+ ( '.' [0-9]+ )? [smh] )+ ;
NUMBER: '0x'?[0-9]+ ;
IDENTIFIER: [a-zA-Z_.][a-zA-Z_.0-9]+ ;
GLOB: [-a-zA-Z.0-9*?]+ ;
EQUALS: '==' ;
MATCHES: '=~' ;
NOT_EQUALS: '!=';
NOT_MATCHES: '!~';
STRING
 : '\'' ( '\\' . | ~[\\\r\n\f'] )* '\''
 | '"' ( '\\' . | ~[\\\r\n\f"] )* '"'
 ;
REGEXP
 : '/' ( '\\' . | ~[\\\r\n\f/] )* '/'
 ;

fragment COMMENT: '#' ~('\n')+;
fragment SPACES: ' '+ ;

SKIP_ : ( SPACES | COMMENT ) -> skip ;

prog : line* EOF ;
line : ( run | set | add | remove | list )? '\n' ;
run : RUN ;
set: SET varname=IDENTIFIER EQUALS? varvalue=value ;
add: ADD HOSTS glob=GLOB filters=filter* ;
remove: REMOVE HOSTS glob=GLOB filters=filter* ;
list: LIST HOSTS oneline=ONELINE? ;
filter: IDENTIFIER ( ( EQUALS | NOT_EQUALS ) value | ( MATCHES | NOT_MATCHES ) REGEXP );
value: NUMBER | STRING | DURATION | IDENTIFIER ;
