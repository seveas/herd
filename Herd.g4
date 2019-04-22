grammar Herd;

// Keep this  the top, as it gobbles up everything after it
RUN: 'run' ~('\n')*;
SET: 'set' ;
ADD: 'add' ;
HOSTS: 'hosts' ;
DURATION: ( '-'? [0-9]+ ( '.' [0-9]+ )? [smh] )+ ;
NUMBER: [0-9]+ ;
IDENTIFIER: [a-zA-Z_.][a-zA-Z_.0-9]+ ;
GLOB: [-a-zA-Z.0-9*]+ ;
EQUALS: '=' ;
STRING
 : '\'' ( '\\' . | ~[\\\r\n\f'] )* '\''
 | '"' ( '\\' . | ~[\\\r\n\f"] )* '"'
 ;

fragment COMMENT: '#' ~('\n')+;
fragment SPACES: ' '+ ;

SKIP_ : ( SPACES | COMMENT ) -> skip ;

prog : line* EOF ;
line : ( run | set | add )? '\n' ;
run : RUN ;
set: SET varname=IDENTIFIER EQUALS? varvalue=value ;
add: ADD HOSTS glob=GLOB filters=filter* ;
filter: IDENTIFIER EQUALS value ;
value: NUMBER | STRING | DURATION ;
